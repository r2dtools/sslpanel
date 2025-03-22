import api, { configWithAuth, getErrorMessage } from '../../lib/api';
import { CertificatesRequest, CertificateMap, DownloadCertificateRequest, DownloadCertificateResponse } from './types';

export const getCertificatesApi = async (request: CertificatesRequest) => {
    try {
        const response = await api.get(`/v1/modules/certificates/${request.guid}/storage/certificates`, configWithAuth(request.token));

        return response.data.certificates as CertificateMap;
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const downloadCertificateApi = async (request: DownloadCertificateRequest) => {
    try {
        const data = {
            name: request.name,
            token: request.token,
        };
        const response = await api.post(`/v1/modules/certificates/${request.guid}/storage/download`, data, configWithAuth(request.token));

        return response.data as DownloadCertificateResponse;
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};
