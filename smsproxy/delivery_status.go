package smsproxy

type MessageStatus string

var NotFound = MessageStatus("NOT FOUND")
var Accepted = MessageStatus("ACCEPTED")
var Confirmed = MessageStatus("CONFIRMED")
var Failed = MessageStatus("FAILED")
var Delivered = MessageStatus("DELIVERED")

var finalStatuses = []MessageStatus{Failed, Delivered}
var allStatuses = []MessageStatus{NotFound, Accepted, Confirmed, Failed, Delivered}
