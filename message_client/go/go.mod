module wraith.me/clientside_crypto

go 1.21.0

toolchain go1.21.1

require wraith.me/message_server v0.0.0

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/norunners/vert v0.0.0-20221203075838-106a353d42dd // indirect
	golang.org/x/crypto v0.25.0 // indirect
)

replace wraith.me/message_server v0.0.0 => ../../message_server
