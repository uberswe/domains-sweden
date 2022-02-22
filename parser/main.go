package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/uberswe/domains-sweden/models"
	"github.com/uberswe/domains-sweden/queue"
	"gorm.io/gorm"
	"log"
	"time"
)

type Event struct {
	URL  string
	Type string
	Time time.Time
}

type Parser struct {
	DB    *gorm.DB
	Queue *queue.Queue
}

func New(db *gorm.DB) *Parser {
	q := queue.NewQueue("domain_parser", 10000000)
	p := Parser{
		DB:    db,
		Queue: q,
	}
	go p.run()
	go p.hearthbeat()
	return &p
}

func (p *Parser) run() {
	w := queue.NewWorker(p.Queue)
	if !w.DoWork() {
		log.Println("Finished running domain parser queue")
	}
}

func (p *Parser) hearthbeat() {
	for range time.Tick(time.Minute * 10) {
		var domains []models.Domain
		res := p.DB.Model(models.Domain{}).Order("RAND()").Limit(10).Find(&domains)
		if res.Error != nil {
			log.Println(res.Error)
		} else {
			for _, d := range domains {
				p.Parse(d)
			}
		}
	}
}

func (p *Parser) Parse(d models.Domain) {
	var payload []byte
	payload, err := json.Marshal(d)
	if err != nil {
		log.Println(err)
		return
	}
	j := queue.Job{
		Name:    d.Host,
		Payload: payload,
		Action:  p.jobTrigger,
	}
	p.Queue.AddJob(j)
}

func (p *Parser) jobTrigger(payload []byte) error {
	var d models.Domain
	err := json.Unmarshal(payload, &d)
	if err != nil {
		return err
	}

	parse := models.Parse{
		DomainID: d.ID,
	}
	res := p.DB.Where(parse).First(&parse)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return res.Error
	}

	// We only parse sites every 14 days max
	if time.Now().Before(parse.UpdatedAt.Add(14 * 24 * time.Hour)) {
		return errors.New("site has been parsed in the last 14 days")
	}

	url := fmt.Sprintf("https://%s", d.Host)
	content, requestSize, responseTime, screenshot, blurredScreenshot, events, requested, err2 := p.process(url)
	if err2 != nil {
		errString := err2.Error()
		parse.Requested = requested
		parse.Error = &errString

		res = p.DB.Save(&parse)
		if res.Error != nil {
			return res.Error
		}
		return err2
	}

	parse.Content = []byte(content)
	parse.Size = requestSize
	parse.LoadTime = responseTime
	parse.Screenshot = screenshot
	parse.BlurredScreenshot = blurredScreenshot
	parse.Requested = requested
	for _, e := range events {
		parse.Events = append(parse.Events, models.ParseEvent{
			URL:  e.URL,
			Type: e.Type,
			Time: e.Time,
		})
	}

	res = p.DB.Save(&parse)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
