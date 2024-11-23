import api, { configWithAuth, getErrorMessage } from '../../lib/api';
import { SaveServerRequest, Server } from './types';

export const getServersApi = async (token: string) => {
    try {
        const response = await api.get('/v1/servers', configWithAuth(token));

        return response.data.servers as Server[];
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const addServerApi = async (server: SaveServerRequest, token: string) => {
    try {
        await api.post('/v1/servers', server, configWithAuth(token));
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const editServerApi = async (id: number, server: SaveServerRequest, token: string) => {
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
