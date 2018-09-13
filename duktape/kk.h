
struct kk_ptr {
	void * ptr;
};

struct kk_ptr * kk_push_ptr(struct duk_hthread *ctx);
struct kk_ptr * kk_to_ptr(struct duk_hthread *ctx,duk_idx_t idx);

