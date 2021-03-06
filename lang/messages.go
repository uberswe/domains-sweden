package lang

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var translationMessages = []i18n.Message{
	{
		ID:    "site_name",
		Other: "Domäner.xyz",
	},
	{
		ID:    "home",
		Other: "Home",
	},
	{
		ID:    "activation_validation_token",
		Other: "Please provide a valid activation token",
	},
	{
		ID:    "activation_success",
		Other: "Account activated. You may now proceed to login to your account.",
	},
	{
		ID:    "activate",
		Other: "Activate",
	},
	{
		ID:    "admin",
		Other: "Admin",
	},
	{
		ID:    "forgot_password",
		Other: "Forgot Password",
	},
	{
		ID:    "forgot_password_success",
		Other: "An email with instructions describing how to reset your password has been sent.",
	},
	{
		ID:    "password_reset",
		Other: "Password Reset",
	},
	{
		ID:    "password_reset_email",
		Other: "Use the following link to reset your password. If this was not requested by you, please ignore this email.\n%s",
	},
	{
		ID:    "login",
		Other: "Login",
	},
	{
		ID:    "login_error",
		Other: "Could not login, please make sure that you have typed in the correct email and password. If you have forgotten your password, please click the forgot password link below.",
	},
	{
		ID:    "login_activated_error",
		Other: "Account is not activated yet.",
	},
	{
		ID:    "404_not_found",
		Other: "404 Not Found",
	},
	{
		ID:    "register",
		Other: "Register",
	},
	{
		ID:    "password_error",
		Other: "Your password must be 8 characters in length or longer",
	},
	{
		ID:    "register_error",
		Other: "Could not register, please make sure the details you have provided are correct and that you do not already have an existing account.",
	},
	{
		ID:    "register_success",
		Other: "Thank you for registering. An activation email has been sent with steps describing how to activate your account.",
	},
	{
		ID:    "user_activation",
		Other: "User Activation",
	},
	{
		ID:    "user_activation_email",
		Other: "Use the following link to activate your account. If this was not requested by you, please ignore this email.\n%s",
	},
	{
		ID:    "resend_activation_email_subject",
		Other: "Resend Activation Email",
	},
	{
		ID:    "resend_activation_email_success",
		Other: "A new activation email has been sent if the account exists and is not already activated. Please remember to check your spam inbox in case the email is not showing in your inbox.",
	},
	{
		ID:    "reset_password",
		Other: "Reset Password",
	},
	{
		ID:    "reset_password_error",
		Other: "Could not reset password, please try again",
	},
	{
		ID:    "password_reset_success",
		Other: "Your password has successfully been reset.",
	},
	{
		ID:    "search",
		Other: "Search",
	},
	{
		ID:    "search_results",
		Other: "Search Results",
	},
	{
		ID:    "no_results_found",
		Other: "No results found",
	},
	{
		ID:    "404_message_1",
		Other: "The page you're looking for could not be found.",
	},
	{
		ID:    "click_here",
		Other: "Click here",
	},
	{
		ID:    "404_message_2",
		Other: "to return to the main page.",
	},
	{
		ID:    "admin_dashboard",
		Other: "Admin Dashboard",
	},
	{
		ID:    "dashboard_message",
		Other: "You now have an authenticated session, feel free to log out using the link in the navbar above.",
	},
	{
		ID:    "footer_message_1",
		Other: "Fork this project on",
	},
	{
		ID:    "created_by",
		Other: "Created by",
	},
	{
		ID:    "forgot_password",
		Other: "Forgot password?",
	},
	{
		ID:    "forgot_password_message",
		Other: "Use the form below to reset your password. If we have an account with your email you will receive instructions on how to reset your password.",
	},
	{
		ID:    "email_address",
		Other: "Email address",
	},
	{
		ID:    "request_reset_email",
		Other: "Request reset email",
	},
	{
		ID:    "lang_key",
		Other: "en",
	},
	{
		ID:    "home",
		Other: "Home",
	},
	{
		ID:    "admin",
		Other: "Admin",
	},
	{
		ID:    "logout",
		Other: "Logout",
	},
	{
		ID:    "login",
		Other: "Login",
	},
	{
		ID:    "register",
		Other: "Register",
	},
	{
		ID:    "search",
		Other: "Search",
	},
	{
		ID:    "index_message_1",
		Other: "A simple website with user login and registration.",
	},
	{
		ID:    "index_message_2",
		Other: "The frontend uses",
	},
	{
		ID:    "index_message_3",
		Other: "and the backend is written in",
	},
	{
		ID:    "index_message_4",
		Other: "Read more about this project on",
	},
	{
		ID:    "password",
		Other: "Password",
	},
	{
		ID:    "login_terms",
		Other: "By pressing the button below to login you agree to the use of cookies on this website.",
	},
	{
		ID:    "request_new_activation_email",
		Other: "Request a new activation email",
	},
	{
		ID:    "resend_activation_email",
		Other: "Resend Activation Email",
	},
	{
		ID:    "resend_activation_email_message",
		Other: "If you have already registered but never activated your account you can use the form below to request a new activation email.",
	},
	{
		ID:    "request_activation_email",
		Other: "Request activation email",
	},
	{
		ID:    "reset_password_message",
		Other: "Please enter a new password.",
	},
	{
		ID:    "domain_first_seen",
		Other: "This domain was first seen by Domäner.xyz at",
	},
	{
		ID:    "domain_release_at_1",
		Other: "This domain is pending delete and if it is not renewed it can be registered again after",
	},
	{
		ID:    "domain_release_at_2",
		Other: "at 04.00 UTC at the earliest",
	},
	{
		ID:    "nameservers",
		Other: "Nameservers",
	},
	{
		ID:    "expired_message",
		Other: "Expired Domains that may become available soon",
	},
	{
		ID:    "domain",
		Other: "Domain",
	},
	{
		ID:    "domains",
		Other: "Domains",
	},
	{
		ID:    "earliest_availability_date",
		Other: "Earliest Availability Date",
	},
	{
		ID:    "nameserver_message",
		Other: "This nameserver was first seen by Domäner.xyz at",
	},
	{
		ID:    "previous",
		Other: "Previous",
	},
	{
		ID:    "next",
		Other: "Next",
	},
	{
		ID:    "nameserver",
		Other: "Nameserver",
	},
	{
		ID:    "most_popular_nameservers",
		Other: "The most popular nameservers",
	},
	{
		ID:    "search_all_domains",
		Other: "Search All Swedish .SE And .NU Domains",
	},
	{
		ID:    "see_all_release_domains",
		Other: "See all domains being released soon",
	},
	{
		ID:    "see_all_nameservers",
		Other: "See all nameservers",
	},
	{
		ID:    "tools",
		Other: "Tools",
	},
	{
		ID:    "data",
		Other: "Data",
	},
	{
		ID:    "language",
		Other: "Language",
	},
	{
		ID:    "whois_lookup",
		Other: "WHOIS Lookup",
	},
	{
		ID:    "top_nameservers",
		Other: "Top Nameservers",
	},
	{
		ID:    "domains_being_released_soon",
		Other: "Domains Being Released Soon",
	},
	{
		ID:    "min_length",
		Other: "Min Length",
	},
	{
		ID:    "max_length",
		Other: "Max Length",
	},
	{
		ID:    "extension",
		Other: "Extension",
	},
	{
		ID:    "website",
		Other: "Website",
	},
	{
		ID:    "days_to_release",
		Other: "Days to release",
	},
	{
		ID:    "releasing_soon",
		Other: "Releasing soon",
	},
	{
		ID:    "no_special_characters",
		Other: "No special characters",
	},
	{
		ID:    "no_numbers",
		Other: "No numbers",
	},
	{
		ID:    "request_whois_lookup",
		Other: "Request WHOIS Lookup",
	},
	{
		ID:    "visit",
		Other: "Visit",
	},
	{
		ID:    "website_data",
		Other: "Website Data",
	},
	{
		ID:    "page_load_time",
		Other: "Page Load Time",
	},
	{
		ID:    "seconds",
		Other: "Seconds",
	},
	{
		ID:    "page_size",
		Other: "Page Size",
	},
	{
		ID:    "mb",
		Other: "Mb",
	},
	{
		ID:    "last_retrieved",
		Other: "Last Retrieved",
	},
	{
		ID:    "unknown",
		Other: "Unknown",
	},
	{
		ID:    "ok_caps",
		Other: "OK",
	},
	{
		ID:    "error",
		Other: "Error",
	},
	{
		ID:    "website_status",
		Other: "Website Status",
	},
	{
		ID:    "has_had_a_website",
		Other: "Has had a website",
	},
	{
		ID:    "has_never_had_a_website",
		Other: "Has never had a website",
	},
}
