export interface Server {
    id: number;
    guid: string;
    name: string;
    os_code: string;
    os_version: string;
    ipv4_address: string;
    ipv6_address: string;
    agent_version: string;
    agent_port: number;
    is_active: number;
    is_registered: number;
    account_id: number;
    created_at: string;
};

export enum OsCode {
    Ubuntu = 'ubuntu',
    Censos = 'centos',
    Debian = 'debian',
}

export interface SaveServerRequest {
    name: string;
    ipv4_address: string;
    ipv6_address: string;
    agent_port: number;
}

export interface ServerSavePayload extends SaveServerRequest {
    token: string;
    id?: number;
}

export interface ServerDeletePayload {
    id: number;
    token: string;
}
