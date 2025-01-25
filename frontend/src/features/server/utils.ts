import centosIcon from '../../images/server/centos-icon.svg';
import debianIcon from '../../images/server/debian-icon.svg';
import ubuntuIcon from '../../images/server/ubuntu-icon.svg';
import linuxIcon from '../../images/server/linux-icon.svg';
import apacheIcon from '../../images/webserver/apache.svg';
import nginxIcon from '../../images/webserver/nginx.svg';
import defaultWebServerIcon from '../../images/webserver/default.png';
import { DomainCertificate, OsCode } from './types';
// @ts-ignore
import * as caIcons from '../../images/ca';
import moment from 'moment';
import {
    AMAZON_CODE,
    APACHE_CODE,
    CLOUD_FLARE_CODE,
    COMODO_CODE,
    DIGICERT_CODE,
    LE_CODE,
    NGINX_CODE,
    SECTIGO_CODE,
} from './constants';

export const getOsIcon = (code: string): string => {
    let icon = '';

    switch (code) {
        case OsCode.Ubuntu:
            icon = ubuntuIcon;

            break;
        case OsCode.Censos:
            icon = centosIcon;

            break;
        case OsCode.Debian:
            icon = debianIcon;

            break;
        default:
            icon = linuxIcon;
    }

    return icon;
};

export const getOsName = (code: string): string => {
    let name = '';

    switch (code) {
        case OsCode.Ubuntu:
            name = 'Ubuntu';

            break;
        case OsCode.Censos:
            name = 'CentOs';

            break;
        case OsCode.Debian:
            name = 'Debian';

            break;
        default:
            name = 'Linux';
    }

    return name;
};

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
    }

    return null;
};

const checkOrganization = (organization: string[], name: string) => {
    return !!organization.find(org => org.indexOf(name) !== -1);
};

export const getWebServerIcon = (code?: string | null): string => {
    if (code === APACHE_CODE) {
        return apacheIcon;
    } else if (code === NGINX_CODE) {
        return nginxIcon;
    }

    return defaultWebServerIcon;
};

export const getCertificateIssuerIcon = (code?: string | null): string | null => {
    if (!code) {
        return null;
    }

    return caIcons[code] || null;
};
