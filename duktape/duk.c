
#include "duk_config.h"
#include "duktape.h"
#include "duk.h"

struct Context * NewContext() {

    struct Context * v = malloc(sizeof(struct Context));
    v->ctx = duk_create_heap_default();
    v->scope = NULL;

    duk_push_global_object(v->ctx);
    duk_push_string(v->ctx,"__Context");
    duk_push_pointer(v->ctx,v);
    duk_put_prop(v->ctx,-3);
    duk_pop(v->ctx);

    return v;
}

void RecycleContext(struct Context * v) {
    duk_destroy_heap(v->ctx);
    free(v);
}

