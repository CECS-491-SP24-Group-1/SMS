const TEST_EP = "https://webhook.site/f3479074-d1b5-465a-afcb-a58bdc884dfd";
const LOCAL_EP_BASE = "http://localhost:8888";

const WASM_URL = "./static/wasm/ed25519_keygen.wasm";
const WASM_PROD_URL = "./static/wasm/ed25519_keygen.min.wasm";

const REGISTER_EP = `${LOCAL_EP_BASE}/auth/register`;

const AUTH_TEST_EP = `${LOCAL_EP_BASE}/auth/test`;

const LOGIN_S1_EP = `${LOCAL_EP_BASE}/auth/login_req`;
const LOGIN_S2_EP = `${LOCAL_EP_BASE}/auth/login_verify`;

const REFRESH_EP = `${LOCAL_EP_BASE}/auth/refresh`;