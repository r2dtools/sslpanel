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

export interface DomainSettings {
    commondirstatus: CommonDirStatus;
    renewal: boolean;
}

export interface DomainSecureRequest {
    guid: string;
    email: string;
    subjects: string[];
    servername: string;
    webserver: string;
    challengetype: string;
    assign: boolean;
    token: string;
};

export interface DomainSecurePayload {
    guid: string;
    email: string;
    subjects: string[];
    servername: string;
    webserver: string;
    challengetype: string;
    assign: boolean;
    token: string;
};

export interface DomainFetchPayload {
    guid: string;
    domainname: string;
    token: string;
}

export interface CommonDirStatusRequest {
    guid: string;
    webserver: string;
    servername: string;
    token: string;
}

export interface DomainSettingsPayload {
    guid: string;
    domain: Domain;
    token: string;
}

export interface CommonDirStatus {
    status: boolean;
}

export interface ChangeCommonChallengeDirStatusRequest {
    guid: string;
    webserver: string;
    servername: string;
    status: boolean;
    token: string;
}

export interface ChangeCommonChallengeDirStatusPayload {
    guid: string;
    domain: Domain;
    token: string;
    status: boolean;
}

export interface DomainConfigRequest {
    guid: string;
    webserver: string;
    domainname: string;
    token: string;
}

export interface DomainConfigPayload {
    guid: string;
    token: string;
    domain: Domain;
}

export interface DomainRequest {
    guid: string;
    domainname: string;
    webserver?: string;
    token: string;
}
