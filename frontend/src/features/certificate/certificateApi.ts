import api, { configWithAuth, getErrorMessage } from '../../lib/api';
import {
    CertificatesRequest,
    CertificateMap,
    DownloadCertificateRequest,
    DownloadCertificateResponse,
    UploadCertificateRequest,
    GenerateSelfSignedCertificateRequest,
} from './types';

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
        };
        const response = await api.post(`/v1/modules/certificates/${request.guid}/storage/download`, data, configWithAuth(request.token));

        return response.data as DownloadCertificateResponse;
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const uploadCertificateApi = async (request: UploadCertificateRequest) => {
    try {
        const formData = new FormData();
        formData.append('name', request.name);
        formData.append('file', request.file);

        await api.post(`/v1/modules/certificates/${request.guid}/storage/upload`, formData, configWithAuth(request.token, {
            headers: {
                'Content-Type': 'multipart/form-data',
            },
        }));
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const generateSelfSignedCertificateApi = async (request: GenerateSelfSignedCertificateRequest) => {
    try {
        const data = {
            certName: request.certName,
            commonName: request.commonName,
            email: request.email,
            country: request.country,
            province: request.province,
            locality: request.locality,
            altNames: request.altNames,
            organization: request.organization,
        };

        await api.post(`/v1/modules/certificates/${request.guid}/storage/add-self-signed`, data, configWithAuth(request.token));
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};
