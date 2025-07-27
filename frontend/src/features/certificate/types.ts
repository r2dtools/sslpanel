export interface StorageCertificateItem {
    name: string;
    storage: string;
    certificate: Certificate;
}

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

export interface CertificatesRequest {
    guid: string;
    token: string;
}

export interface CertificatesPayload {
    guid: string;
    token: string;
}

export interface DownloadCertificateRequest {
    guid: string;
    token: string;
    name: string;
    storage: string;
}

export interface DownloadCertificateResponse {
    name: string;
    content: string;
}

export interface UploadCertificateRequest {
    name: string;
    file: File;
    token: string;
    guid: string;
}

export interface UploadCertificatePayload {
    name: string;
    file: File;
    token: string;
    guid: string;
}

export interface SelfSignedCertificateFormData {
    certName: string,
    commonName: string,
    email: string | null,
    country: string | null,
    province: string | null,
    locality: string | null,
    altNames: string[],
    organization: string | null,
}

export interface GenerateSelfSignedCertificateRequest extends SelfSignedCertificateFormData {
    token: string;
    guid: string;
}

export interface GenerateSelfSignedCertificatePayload extends SelfSignedCertificateFormData {
    token: string;
    guid: string;
}
