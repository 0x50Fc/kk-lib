#include "duk_config.h"
#include "duktape.h"
#include "kk.h"

struct kk_ptr * kk_push_ptr(struct duk_hthread *ctx) {
	return (struct kk_ptr *) duk_push_fixed_buffer(ctx,sizeof(struct kk_ptr));
}

struct kk_ptr * kk_to_ptr(struct duk_hthread *ctx,duk_idx_t idx) {
	size_t n=0;
	return (struct kk_ptr *) duk_to_buffer(ctx,idx,&n);
}

