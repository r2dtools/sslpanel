import api, { configWithAuth, getErrorMessage } from '../../lib/api';
import { ServerSaveRequest, Server, ServerDetails, ChangeSettingRequest, RenewalLog } from './types';

export const getServersApi = async (token: string) => {
    try {
        const response = await api.get('/v1/servers', configWithAuth(token));

        return response.data.servers as Server[];
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const getServerApi = async (guid: string, token: string) => {
    try {
        const response = await api.get(`/v1/servers/${guid}`, configWithAuth(token));

        return response.data.server as Server;
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const getServerDetailsApi = async (guid: string, token: string) => {
    try {
        const response = await api.get(`/v1/servers/${guid}/details`, configWithAuth(token));

        return response.data.server as ServerDetails;
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const addServerApi = async (server: ServerSaveRequest, token: string) => {
    try {
        await api.post('/v1/servers', server, configWithAuth(token));
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const editServerApi = async (id: number, server: ServerSaveRequest, token: string) => {
    try {
        await api.post(`/v1/servers/${id}`, server, configWithAuth(token));
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const deleteServerApi = async (id: number, token: string) => {
    try {
        await api.delete(`/v1/servers/${id}`, configWithAuth(token));
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const changeCertbotStatusApi = async (request: ChangeSettingRequest) => {
    try {
        const data = {
            value: request.value === 'true',
        };
        const response = await api.post(`/v1/servers/${request.guid}/settings/certbot-status`, data, configWithAuth(request.token));

        return response.data.version as string
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const getServerRenewalLogsApi = async (guid: string, token: string) => {
    try {
        const response = await api.get(`/v1/modules/certificates/${guid}/renewal/latest-logs`, configWithAuth(token));

        return response.data.logs as RenewalLog[];
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};
