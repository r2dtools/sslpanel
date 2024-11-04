import api, { getErrorMessage } from "../../lib/api";
import { LoginResponse, User } from "./types";

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
        const response = await api.post("/v1/auth/me", null, {
            headers: {
                Authorization: 'Bearer ' + token,
            }
        });

        return response.data as User;
    } catch (error) {
        throw new Error(`Authorization failed: ${getErrorMessage(error)}`)
    }
};
