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
    account_id: number;
    created_at: string;
    token: string;
};

export interface ServerDetails extends Server {
    hostname: string;
    os: string;
    platform_family: string;
    kernal_version: string;
    kernal_arch: string;
    virtualization: string;
    uptime: number;
    boottime: number;
    domains: Domain[];
}

export enum OsCode {
    Ubuntu = 'ubuntu',
    Censos = 'centos',
    Debian = 'debian',
}

export interface ServerSaveRequest {
    name: string;
    ipv4_address: string;
    ipv6_address: string;
    agent_port: number;
    token: string;
}

export interface ServerSavePayload extends ServerSaveRequest {
    authToken: string;
    id?: number;
    guid?: string;
}

export interface ServerDeletePayload {
    id: number;
    token: string;
}

export interface ServerFetchPayload {
    guid: string;
    token: string;
}

interface DomainAddress {
    isIpv6: boolean;
    host: string;
    port: string;
}

export interface Domain {
    filepath: string;
    servername: string;
    docroot: string;
    webserver: string;
    aliases: string[];
    ssl: boolean;
    addresses: DomainAddress[];
    certificate: DomainCertificate | null;
}

export interface DomainCertificate {
    cn: string;
    validfrom: string;
    validto: string;
    dnsnames: string[] | null;
    emailaddresses: string[] | null;
    organization: string[] | null;
    province: string[] | null;
    country: string[] | null;
    locality: string[] | null;
    isca: boolean;
    isvalid: boolean;
    issuer: Issuer;
}

export interface Issuer {
    cn: string;
    organization: string[] | null;
}
