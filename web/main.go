package main

import (
	"nobozo/config"
	"nobozo/handle"
)

const (
	f_address1 = iota
	f_address2
	f_authtoken
	f_city
	f_comments
	f_country
	f_created
	f_creator
	f_emailaddress
	f_freeformcontactinfo
	f_gecos
	f_homephone
	f_id
	f_lang
	f_lastupdated
	f_mobilephone
	f_name
	f_nickname
	f_organization
	f_pagerphone
	f_password
	f_realname
	f_signature
	f_smimecertificate
	f_state
	f_timezone
	f_workphone
	f_zip
)

type currentUser_t struct {
	sbRecordPrimaryRecordCacheKey string
	quotedTable                   string
}

type systemUser_t struct {
	currentUser *currentUser_t
	fetched     fetched_t
	values      values_t
	table       string
}

type fetched_t [f_zip - 1]bool
type values_t [f_zip - 1]string

func main() {
	connectToDatabase()
	initSystemObjects()
	config.LoadConfigFromDatabase()
}

func connectToDatabase() {
	handle.Connect()
}

func initSystemObjects() {
}
