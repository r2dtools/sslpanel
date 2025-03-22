import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'
import { RootState } from '../../app/store';
import { FetchStatus } from '../../app/types';
import { toast } from 'react-toastify';
import {
    CertificateMap,
    CertificatesPayload,
    GenerateSelfSignedCertificatePayload,
    UploadCertificatePayload,
} from './types';
import { generateSelfSignedCertificateApi, getCertificatesApi, uploadCertificateApi } from './certificateApi';

export interface CertificatesState {
    certificates: CertificateMap;
    certificatesStatus: FetchStatus;
    certificateUploadStatus: FetchStatus;
    selfSignedCertificateGenerateStatus: FetchStatus;
}

const initialState: CertificatesState = {
    certificates: {},
    certificatesStatus: FetchStatus.Idle,
    certificateUploadStatus: FetchStatus.Idle,
    selfSignedCertificateGenerateStatus: FetchStatus.Idle,
}

export const fetchCertificates = createAsyncThunk(
    'certificate/list',
    async (payload: CertificatesPayload) => {
        return await getCertificatesApi({
            guid: payload.guid,
            token: payload.token,
        });
    },
);

export const uploadCertificate = createAsyncThunk(
    'certificate/upload',
    async (payload: UploadCertificatePayload) => {
        await uploadCertificateApi({
            guid: payload.guid,
            token: payload.token,
            name: payload.name,
            file: payload.file,
        });
    },
);

export const generateSelfSignedCertificate = createAsyncThunk(
    'certificate/add-self-signed',
    async (payload: GenerateSelfSignedCertificatePayload) => {
        await generateSelfSignedCertificateApi({
            guid: payload.guid,
            token: payload.token,
            certName: payload.certName,
            commonName: payload.commonName,
            email: payload.email,
            country: payload.country,
            province: payload.province,
            locality: payload.locality,
            organization: payload.organization,
            altNames: payload.altNames,
        });
    },
);

export const certificatesSlice = createSlice({
    name: 'certificates',
    initialState,
    reducers: {},
    extraReducers: builder => {
        builder
            .addCase(fetchCertificates.pending, state => {
                state.certificatesStatus = FetchStatus.Pending;
            })
            .addCase(fetchCertificates.fulfilled, (state, action) => {
                state.certificatesStatus = FetchStatus.Succeeded;
                state.certificates = action.payload;
            })
            .addCase(fetchCertificates.rejected, (state, action) => {
                state.certificatesStatus = FetchStatus.Failed;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            }).addCase(uploadCertificate.pending, state => {
                state.certificateUploadStatus = FetchStatus.Pending;
            })
            .addCase(uploadCertificate.fulfilled, (state) => {
                state.certificateUploadStatus = FetchStatus.Succeeded;
            })
            .addCase(uploadCertificate.rejected, (state, action) => {
                state.certificateUploadStatus = FetchStatus.Failed;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            }).addCase(generateSelfSignedCertificate.pending, state => {
                state.selfSignedCertificateGenerateStatus = FetchStatus.Pending;
            })
            .addCase(generateSelfSignedCertificate.fulfilled, (state) => {
                state.selfSignedCertificateGenerateStatus = FetchStatus.Succeeded;
            })
            .addCase(generateSelfSignedCertificate.rejected, (state, action) => {
                state.selfSignedCertificateGenerateStatus = FetchStatus.Failed;

                if (action.error.message) {
                    toast.error(action.error.message);
                }
            });
    },
});

export const selectCertificates = (state: RootState) => state.certificates.certificates;
export const selectCertificatesFetchStatus = (state: RootState) => state.certificates.certificatesStatus;
export const selectCertificateUplaodStatus = (state: RootState) => state.certificates.certificateUploadStatus;
export const selectSelfSignedCertificateGenerateStatus = (state: RootState) => state.certificates.selfSignedCertificateGenerateStatus;

export default certificatesSlice.reducer
