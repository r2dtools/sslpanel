import PasswordValidator from 'password-validator';

export const PASSWORD_HELPER_TEXT = 'The password must be at least 8 characters long, must not contain spaces, contain at least 1 uppercase and lowercase letter, and contain at least 2 numbers.'

export const createPasswordValidator = (): PasswordValidator => {
    const validator = new PasswordValidator();

    return validator
        .is().min(8)
        .is().max(100)
        .has().uppercase()
        .has().lowercase()
        .has().digits(2)
        .has().not().spaces()
};
