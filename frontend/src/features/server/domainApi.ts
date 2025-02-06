import api, { configWithAuth, getErrorMessage } from '../../lib/api';
import {
    CommonChallengeDirStatusChangeRequest,
    CommonDirStatus,
    DomainCertificate,
    DomainSecureRequest,
    GetCommonDirStatusRequest,
    GetDomainConfigRequest,
} from './types';

export const secureDomainApi = async (guid: string, certData: DomainSecureRequest, token: string) => {
    try {
        const response = await api.post(`/v1/modules/certificates/${guid}/domain/issue`, certData, configWithAuth(token));

        return response.data.certificate as DomainCertificate;
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const getCommonDirApi = async (guid: string, requestData: GetCommonDirStatusRequest, token: string) => {
    try {
        const response = await api.post(`/v1/modules/certificates/${guid}/domain/commondir-status`, requestData, configWithAuth(token));

        return response.data.commondir as CommonDirStatus;
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const changeCommonDirStatusApi = async (guid: string, requestData: CommonChallengeDirStatusChangeRequest, token: string) => {
    try {
        await api.post(`/v1/modules/certificates/${guid}/domain/change-commondir-status`, requestData, configWithAuth(token));
    } catch (error) {
        console.log(error);
        throw new Error(getErrorMessage(error))
    }
};

export const getDomainConfigApi = async (guid: string, requestData: GetDomainConfigRequest, token: string) => {
    try {
        const response = await api.get(`/v1/servers/${guid}/domain-config`, {
            ...configWithAuth(token),
            params: requestData,
        });

        return response.data.content as string;
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};
