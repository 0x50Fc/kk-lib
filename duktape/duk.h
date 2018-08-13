
struct Context {
	duk_context * ctx;
	void * scope;
};

struct Context * NewContext();

void RecycleContext(struct Context * v);
