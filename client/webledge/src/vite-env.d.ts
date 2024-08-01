/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_APP_TITLE: string
    readonly VITE_APP_MIMIRLEDGER_API_URL: string
    // more env variables...
}

interface ImportMeta {
    readonly env: ImportMetaEnvw
}