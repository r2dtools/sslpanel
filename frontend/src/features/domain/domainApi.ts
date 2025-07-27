import { Base64 } from 'js-base64';
import api, { configWithAuth, getErrorMessage } from '../../lib/api';
import {
    DomainSecureRequest,
    ChangeSettingRequest,
    CommonDirStatus,
    Domain,
    DomainCertificate,
    CommonDirStatusRequest,
    DomainConfigRequest,
    DomainRequest,
    DomainSettingsRequest,
    DomainSettingsResponse,
    AssignCertificateRequest,
} from './types';

export const issueCertificateApi = async (request: DomainSecureRequest) => {
    try {
        const domainname = Base64.encodeURI(request.servername);
        const data = {
            email: request.email,
            subjects: request.subjects,
            webserver: request.webserver,
            challengetype: request.challengetype,
            assign: request.assign,
        };
        const response = await api.post(`/v1/modules/certificates/${request.guid}/domain/${domainname}/issue`, data, configWithAuth(request.token));

        return response.data.certificate as DomainCertificate;
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const assignCertificateApi = async (request: AssignCertificateRequest) => {
    try {
        const domainname = Base64.encodeURI(request.servername);
        const data = {
            name: request.name,
            storage: request.storage,
            webserver: request.webserver,
        };
        const response = await api.post(`/v1/modules/certificates/${request.guid}/domain/${domainname}/assign`, data, configWithAuth(request.token));

        return response.data.certificate as DomainCertificate;
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const getCommonDirApi = async (request: CommonDirStatusRequest) => {
    try {
        const domainname = Base64.encodeURI(request.domainname);
        const response = await api.get(`/v1/modules/certificates/${request.guid}/domain/${domainname}/commondir-status`, {
            ...configWithAuth(request.token),
            params: {
                webserver: request.webserver,
            },
        });

        return response.data.commondir as CommonDirStatus;
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const changeCommonDirStatusApi = async (request: ChangeSettingRequest) => {
    try {
        const domainname = Base64.encodeURI(request.servername);
        const data = {
            webserver: request.webserver,
            status: request.status === 'true',
        };
        await api.post(`/v1/modules/certificates/${request.guid}/domain/${domainname}/commondir-status`, data, configWithAuth(request.token));
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const changeSettingApi = async (request: ChangeSettingRequest) => {
    try {
        const domainname = Base64.encodeURI(request.servername);
        const data = {
            settingname: request.name,
            webserver: request.webserver,
            settingvalue: request.status,
        };
        await api.post(`/v1/servers/${request.guid}/domain/${domainname}/settings`, data, configWithAuth(request.token));
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const getDomainConfigApi = async (request: DomainConfigRequest) => {
    try {
        const domainname = Base64.encodeURI(request.domainname);
        const response = await api.get(`/v1/servers/${request.guid}/domain/${domainname}/config`, {
            ...configWithAuth(request.token),
            params: {
                webserver: request.webserver,
            },
        });

        return response.data.content as string;
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const getDomainApi = async (request: DomainRequest) => {
    try {
        const domainname = Base64.encodeURI(request.domainname);
        const response = await api.get(`/v1/servers/${request.guid}/domain/${domainname}`, configWithAuth(request.token));

        return response.data.domain as Domain;
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const getDomainSettingsApi = async (request: DomainSettingsRequest) => {
    try {
        const domainname = Base64.encodeURI(request.domainname);
        const response = await api.get(`/v1/servers/${request.guid}/domain/${domainname}/settings`, configWithAuth(request.token));

        return response.data.settings as DomainSettingsResponse;
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};
