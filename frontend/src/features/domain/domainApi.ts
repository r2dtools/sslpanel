import { Base64 } from 'js-base64';
import api, { configWithAuth, getErrorMessage } from '../../lib/api';
import {
    DomainSecureRequest,
    ChangeCommonChallengeDirStatusRequest,
    CommonDirStatus,
    Domain,
    DomainCertificate,
    CommonDirStatusRequest,
    DomainConfigRequest,
    DomainRequest,
} from './types';

export const secureDomainApi = async (request: DomainSecureRequest) => {
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

export const getCommonDirApi = async (request: CommonDirStatusRequest) => {
    try {
        const domainname = Base64.encodeURI(request.servername);
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

export const changeCommonDirStatusApi = async (request: ChangeCommonChallengeDirStatusRequest) => {
    try {
        const domainname = Base64.encodeURI(request.servername);
        const data = {
            webserver: request.webserver,
            status: request.status,
        };
        await api.post(`/v1/modules/certificates/${request.guid}/domain/${domainname}/commondir-status`, data, configWithAuth(request.token));
    } catch (error) {
        console.log(error);
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
