ifndef MAGMA_ROOT
MAGMA_ROOT := /home/$(USER)/magma
$(warning MAGMA_ROOT is undefined, using default: $(MAGMA_ROOT))
endif

SWAGGER_V1_SPEC := $(MAGMA_ROOT)/orc8r/cloud/go/services/obsidian/swagger/v1/swagger.yml
export SWAGGER_V1_SPEC

LANGUAGES := go

define add_language
	$(eval LANG := $1)
	GEN_LIST += $(LANG)_gen
endef
_ := $(foreach lang, $(LANGUAGES), $(eval $(call add_language,$(lang))))

gen: $(GEN_LIST)
$(GEN_LIST): %_gen:
	make -C $*/ gen
