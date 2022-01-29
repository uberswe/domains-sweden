# Domains Sweden - Domäner.xyz

[![GoDoc](https://godoc.org/github.com/uberswe/domains-sweden?status.svg)](https://godoc.org/github.com/uberswe/domains-sweden)

This project creates a website which fetches active .se and .nu domains from iis.se.

## Translations

This project uses [go-i18n](https://github.com/nicksnyder/go-i18n) to handle translations.

To update languages first run `goi18n extract` to update `active.en.toml`. Then run `goi18n merge active.*.toml` to generate `translate.*.toml` which can then be translated. Finally, run `goi18n merge active.*.toml translate.*.toml` to merge the translated files into the active files.

## Documentation

See [GoDoc](https://godoc.org/github.com/uberswe/domains-sweden) for further documentation.

## Contributions

Contributions are welcome and greatly appreciated. Please note that I am not looking to add any more features to this project but I am happy to take care of bugfixes, updates and other suggestions. If you have a question or suggestion please feel free to [open an issue](https://github.com/uberswe/domains-sweden/issues/new). To contribute code, please fork this repository, make your changes on a separate branch and then [open a pull request](https://github.com/uberswe/domains-sweden/compare).

For security related issues please see my profile, [@uberswe](https://github.com/uberswe), for ways of contacting me privately. 

## License

Please see the `LICENSE` file in the project repository.