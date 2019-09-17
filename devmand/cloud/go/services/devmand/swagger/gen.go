//go:generate cp $SWAGGER_ROOT/$SWAGGER_COMMON $SWAGGER_COMMON
//go:generate swagger generate model -f swagger.yml -t ../obsidian/ -C $SWAGGER_TEMPLATE
//go:generate rm ./$SWAGGER_COMMON

package swagger
