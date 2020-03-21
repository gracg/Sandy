CC = go


RepositoryName = Sandy
ProjectName = Sandy

gpgKeyID = "57D8B5AC65DD78A8304336CFAFA22EFD5052BA66"

.PHONY: default clean 

default:
	@echo "Building"
	@${CC} build Sandy.go server.go db.go key.go client.go pgp.go
	@${CC} build Installer.go

clean:
	@rm Sandy
	@rm Installer