import api, { configWithAuth, getErrorMessage } from '../../lib/api';
import { DomainCertificate, DomainSecureRequest } from './types';

export const secureDomainApi = async (guid: string, certData: DomainSecureRequest, token: string) => {
    try {
        const response = await api.post(`/v1/modules/certificates/${guid}/domain/issue`, certData, configWithAuth(token));

        return response.data.certificate as DomainCertificate;
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};
