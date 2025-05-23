export interface Account {
    id: number
    created_at: string
};

export interface User {
    id: number
    email: string
    is_active: boolean
    is_account_owner: boolean
    account_id: number
    account: Account
    created_at: string
}

export interface LoginResponse {
    expire: string
    token: string
};
