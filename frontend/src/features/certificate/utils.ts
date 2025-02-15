// @ts-ignore
import * as caIcons from '../../images/ca';
import moment from 'moment';
import {
    AMAZON_CODE,
    CLOUD_FLARE_CODE,
    COMODO_CODE,
    DIGICERT_CODE,
    LE_CODE,
    RAPID_SSL_CODE,
    SECTIGO_CODE,
} from './constants';
import { DomainCertificate } from '../domain/types';

export const getSiteCertExpiredDays = (validTo?: string | null) => {
    if (!validTo) {
        return null;
    }

    const validToDate = moment(validTo);
    const currentDate = moment(new Date());

    return validToDate.diff(currentDate, 'd', true);
};

export const isSelfSignedCertificate = (certificate: DomainCertificate | null) => {
    if (!certificate || !certificate.cn) {
        return false;
    }

    const { issuer } = certificate;

    if (!issuer) {
        return false;
    }

    return certificate.cn === issuer.cn;
}

export const getCertificateIssuerCode = (certificate: DomainCertificate | null) => {
    if (!certificate) {
        return null;
    }

    const { issuer } = certificate;

    if (!issuer) {
        return null;
    }

    const organization = issuer.organization || [];
    const cn = issuer.cn || '';

    if (checkOrganization(organization, 'Let\'s Encrypt') || checkOrganization(organization, 'good guys') || cn.indexOf('Let\'s Encrypt') !== -1 || cn.indexOf('R3') !== -1) {
        return LE_CODE;
    } else if (checkOrganization(organization, 'DigiCert') || cn.indexOf('DigiCert') !== -1) {
        return DIGICERT_CODE;
    } else if (checkOrganization(organization, 'Sectigo') || cn.indexOf('Sectigo') !== -1) {
        return SECTIGO_CODE;
    } else if (checkOrganization(organization, 'Comodo') || cn.indexOf('Comodo') !== -1) {
        return COMODO_CODE;
    } else if (checkOrganization(organization, 'Amazon') || cn.indexOf('Amazon') !== -1) {
        return AMAZON_CODE;
    } else if (checkOrganization(organization, 'Cloudflare') || cn.indexOf('Cloudflare') !== -1) {
        return CLOUD_FLARE_CODE;
    } else if (checkOrganization(organization, 'RapidSSL') || cn.indexOf('RapidSSL') !== -1) {
        return RAPID_SSL_CODE;
    }

    return null;
};

const checkOrganization = (organization: string[], name: string) => {
    return !!organization.find(org => org.indexOf(name) !== -1);
};

export const getCertificateIssuerIcon = (code?: string | null): string | null => {
    if (!code) {
        return null;
    }

    // @ts-ignore
    return caIcons[code] || null;
};
