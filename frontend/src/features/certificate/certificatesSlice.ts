import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'
import { RootState } from '../../app/store';
import { FetchStatus } from '../../app/types';
import { toast } from 'react-toastify';
import { CertificateMap, CertificatesPayload } from './types';
import { getCertificatesApi } from './certificateApi';

export interface CertificatesState {
    certificates: CertificateMap;
    certificatesStatus: FetchStatus;
}

const initialState: CertificatesState = {
    certificates: {},
    certificatesStatus: FetchStatus.Idle,
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
            });
    },
});

export const selectCertificates = (state: RootState) => state.certificates.certificates;
export const selectCertificatesFetchStatus = (state: RootState) => state.certificates.certificatesStatus;

export default certificatesSlice.reducer
