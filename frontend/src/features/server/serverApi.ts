import api, { configWithAuth, getErrorMessage } from '../../lib/api';

export const getServers = async (token: string) => {
    try {
        const response = await api.get('/v1/servers', configWithAuth(token));

        return response.data.servers as Server[];
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};
