import api, { configWithAuth, getErrorMessage } from '../../lib/api';
import { ChangePasswordRequest } from './types';

export const changePasswordApi = async (request: ChangePasswordRequest) => {
    try {
        const data = {
            password: request.password,
            newPassword: request.newPassword,
        };

        return await api.post("/v1/settings/change-password", data, configWithAuth(request.token));
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};
