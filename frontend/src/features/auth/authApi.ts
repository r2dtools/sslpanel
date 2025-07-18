import api, { configWithAuth, getErrorMessage } from '../../lib/api';
import { LoginResponse, User } from './types';

export const login = async (email: string, password: string) => {
    try {
        const response = await api.post("/v1/login", { email, password });

        return response.data as LoginResponse;
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const me = async (token: string) => {
    try {
        const response = await api.post("/v1/auth/me", null, configWithAuth(token));

        return response.data.user as User;
    } catch (error) {
        throw new Error(`Authorization failed: ${getErrorMessage(error)}`)
    }
};

export const signUp = async (email: string, password: string) => {
    try {
        await api.post("/v1/register", { email, password });
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const confirm = async (token: string) => {
    try {
        await api.post("/v1/confirm", { token });
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const recover = async (email: string) => {
    try {
        await api.post("/v1/recover", { email });
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};

export const reset = async (token: string) => {
    try {
        await api.post("/v1/reset", { token });
    } catch (error) {
        throw new Error(getErrorMessage(error))
    }
};