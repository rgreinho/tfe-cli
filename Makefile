TARGET_NAME := tfe-cli
PLATFORMS := linux-amd64 darwin-amd64 windows-amd64
TARGET_NAME_PLATFORM := $(foreach platform,$(PLATFORMS),$(TARGET_NAME)-$(platform))

os_arch_split = $(subst -, ,$@)
os = $(word 3, $(os_arch_split))
arch = $(word 4, $(os_arch_split))

.PHONY: clean
clean:
	rm $(TARGET_NAME_PLATFORM)

release: $(TARGET_NAME_PLATFORM)

$(TARGET_NAME_PLATFORM):
	GOOS=$(os) GOARCH=$(arch) go build -o '$@'
