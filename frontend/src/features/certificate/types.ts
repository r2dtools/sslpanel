export interface Certificate {
    cn: string;
    validfrom: string;
    validto: string;
    dnsnames: string[] | null;
    emails: string[] | null;
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

export interface CertificateMap {
    [key: string]: Certificate,
}

export interface CertificatesRequest {
    guid: string;
    token: string;
}

export interface CertificatesPayload {
    guid: string;
    token: string;
}
