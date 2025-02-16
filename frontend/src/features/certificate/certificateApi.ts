import api, { configWithAuth, getErrorMessage } from '../../lib/api';
import { CertificatesRequest, CertificateMap } from './types';

export const getCertificatesApi = async (request: CertificatesRequest) => {
    try {
        const response = await api.get(`/v1/modules/certificates/${request.guid}/storage/certificates`, configWithAuth(request.token));

        return response.data.certificates as CertificateMap;
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};
