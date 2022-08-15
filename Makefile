IMPORT_PATH = $(shell echo `pwd` | sed "s|^$(GOPATH)/src/||g")
APP_NAME = $(shell echo $(IMPORT_PATH) | sed 's:.*/::')
# for ease of use
CANONICAL_NAME = "reservation_thxx_go"
APP_VERSION = 0.1
MAC_TARGET = ./$(APP_NAME)-mac-$(APP_VERSION)
MAC_EXTERNAL = ./$(APP_NAME)-mac-external-$(APP_VERSION)
LINUX_TARGET = ./$(APP_NAME)-linux-$(APP_VERSION)
LINUX_EXTERNAL = ./$(APP_NAME)-linux-external-$(APP_VERSION)
GO_FILES = $(shell find . -type f -name "*.go")
PID = .pid
#go server port
PORT ?= 9000

build: clean $(LINUX_TARGET) $(LINUX_EXTERNAL)

clean:
	@rm -rf $(MAC_TARGET)
	@rm -rf $(MAC_EXTERNAL)
	@rm -rf $(LINUX_TARGET)
	@rm -rf $(LINUX_EXTERNAL)
	@rm -rf $(APP_NAME)-$(APP_VERSION).zip

$(MAC_TARGET): $(GO_FILES)
	@printf "Building mac go binary ......\n"
	@env GOOS=darwin GOARCH=amd64 go build -o $@

$(MAC_EXTERNAL): $(GO_FILES)
	@printf "Building mac external go binary ......\n"
	@env GOOS=darwin GOARCH=amd64 go build -o $@ ./external

$(LINUX_TARGET): $(GO_FILES)
	@printf "Building linux go binary ......\n"
	@env GOOS=linux GOARCH=amd64 go build -o $@

$(LINUX_EXTERNAL): $(GO_FILES)
	@printf "Building linux external go binary ......\n"
	@env GOOS=linux GOARCH=amd64 go build -o $@ ./external

kill:
	@kill `cat $(PID)` || true

dev: clean $(MAC_TARGET) $(MAC_EXTERNAL) $(LINUX_TARGET) $(LINUX_EXTERNAL) restart
	@printf "\n\nWaiting for the file change\n\n"
	@fswatch --one-per-batch $(GO_FILES) | xargs -n1 -I{} make restart || make kill

restart: kill $(MAC_TARGET)
	@printf "\n\nrestart the app .........\n\n"
	@$(MAC_TARGET) -debug --web=:$(PORT) & echo $$! > $(PID)

dist: clean $(LINUX_TARGET) $(LINUX_EXTERNAL)
	@zip -r -v $(APP_NAME)-$(APP_VERSION).zip $(LINUX_TARGET) $(LINUX_EXTERNAL) \
    	templates static deploy tools

image:
	@printf "\n\nBuilding linux image ......\n\n"
	@docker build -t docker.student.tsinghua.edu.cn/$(CANONICAL_NAME):latest .
