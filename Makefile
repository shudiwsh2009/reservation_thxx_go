IMPORT_PATH = $(shell echo `pwd` | sed "s|^$(GOPATH)/src/||g")
APP_NAME = $(shell echo $(IMPORT_PATH) | sed 's:.*/::')
APP_VERSION = 0.1
MAC_TARGET = ./$(APP_NAME)-mac-$(APP_VERSION)
MAC_EXTERNAL = ./$(APP_NAME)-mac-external-$(APP_VERSION)
LINUX_TARGET = ./$(APP_NAME)-linux-$(APP_VERSION)
LINUX_EXTERNAL = ./$(APP_NAME)-linux-external-$(APP_VERSION)
GO_FILES = $(shell find . -type f -name "*.go")
BUNDLE = public/bundles
ASSETS = $(shell find assets -type f)
PID = .pid
NODE_BIN = $(shell npm bin)
#go server port
PORT ?= 9000
#webpack-dev-server port
DEV_HOT_PORT ?= 8090

build: clean $(BUNDLE) $(MAC_TARGET) $(MAC_EXTERNAL)

clean:
	@rm -rf public/bundles
	@rm -rf $(MAC_TARGET)
	@rm -rf $(MAC_EXTERNAL)
	@rm -rf $(LINUX_TARGET)
	@rm -rf $(LINUX_EXTERNAL)
	@rm -rf $(APP_NAME)-$(APP_VERSION).zip

$(BUNDLE): $(ASSETS)
	@$(NODE_BIN)/webpack --progress --colors

$(MAC_TARGET): $(GO_FILES)
	@printf "Building mac go binary ......\n"
	@env GOOS=darwin GOARCH=amd64 go build -o $@

$(MAC_EXTERNAL): $(GO_FILES)
	@printf "Building mac external go binary ......\n"
	@env GOOS=darwin GOARCH=amd64 go build -race -o $@ ./external

$(LINUX_TARGET): $(GO_FILES)
	@printf "Building linux go binary ......\n"
	@env GOOS=linux GOARCH=amd64 go build -o $@

$(LINUX_EXTERNAL): $(GO_FILES)
	@printf "Building linux external go binary ......\n"
	@env GOOS=linux GOARCH=amd64 go build -o $@ ./external

kill:
	@kill `cat $(PID)` || true

dev: clean $(BUNDLE) restart
	@DEV_HOT=true NODE_ENV=development $(NODE_BIN)/webpack-dev-server --config webpack.config.js &
	@printf "\n\nWaiting for the file change\n\n"
	@fswatch --one-per-batch $(GO_FILES) | xargs -n1 -I{} make restart || make kill

restart: kill $(MAC_TARGET)
	@printf "\n\nrestart the app .........\n\n"
	@$(MAC_TARGET) -debug --web=:$(PORT) --devWeb=:$(DEV_HOT_PORT) & echo $$! > $(PID)

dist: clean $(LINUX_TARGET) $(LINUX_EXTERNAL)
	@NODE_ENV=production $(NODE_BIN)/webpack --progress --colors
	@zip -r -v $(APP_NAME)-$(APP_VERSION).zip $(LINUX_TARGET) $(LINUX_EXTERNAL) \
    	webpack-assets.json public templates static deploy tools
