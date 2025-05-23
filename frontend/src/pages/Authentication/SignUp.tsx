import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { Button, Spinner, TextInput } from 'flowbite-react';
import { createPasswordValidator, PASSWORD_HELPER_TEXT } from '../../lib/form';
import { HiMiniEnvelope, HiMiniLockClosed } from 'react-icons/hi2';
import { signUp } from '../../features/auth/authApi';
import { toast } from 'react-toastify';
import AuthSide from '../../components/AuthSide';

const inputTheme = {
    field: {
        input: {
            colors: {
                gray: 'border-gray-300 bg-gray-50 text-gray-900 placeholder-gray-500 focus:border-blue-500 focus:ring-blue-500 dark:border-form-strokedark dark:bg-form-input dark:text-white dark:placeholder-gray-400 dark:focus:border-blue-500 dark:focus:ring-blue-500',
            },
        },
    }
};

const SignUp: React.FC = () => {
    const [loading, setLoading] = useState(false);
    const [email, setEmail] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [repeatePassword, setRepeatePassword] = useState<string>('');
    const [passwordError, setPasswordError] = useState<string>('');
    const [repeatePasswordError, setRepeatePasswordError] = useState<string>('');
    const [accountCreated, setAccountCreated] = useState<boolean>(false);

    const submitDisabled = loading;

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        const passwordValidator = createPasswordValidator();
        const passwordErrors = passwordValidator.validate(password, { details: true }) as any[];

        if (passwordErrors.length) {
            setPasswordError(passwordErrors[0].message);

            return;
        }

        if (password !== repeatePassword) {
            setRepeatePasswordError('Passwords do not match');

            return;
        }

        try {
            setLoading(true);
            await signUp(email, password);
            setAccountCreated(true);
        } catch (error) {
            resetForm();
            toast.error((error as Error).message);
        } finally {
            setLoading(false);
        }
    };

    const handlePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setPasswordError('');
        setPassword(event.target.value);
    };

    const handleRepeatePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRepeatePasswordError('');
        setRepeatePassword(event.target.value);
    };

    const resetForm = () => {
        setPassword('');
        setRepeatePassword('');
    };

    return (
        <div className="rounded-sm border border-stroke bg-white shadow-default dark:border-strokedark dark:bg-boxdark">
            <div className="flex flex-wrap items-center">
                <div className="hidden w-full xl:block xl:w-1/2">
                    <AuthSide />
                </div>
                <div className="w-full border-stroke dark:border-strokedark xl:w-1/2 xl:border-l-2">
                    <div className="w-full p-4 sm:p-12.5 xl:p-17.5">
                        <h2 className="mb-9 text-2xl font-bold text-black dark:text-white sm:text-title-xl2">
                            Sign Up to SSLPanel
                        </h2>
                        {accountCreated && (
                            <>
                                <div className='font-medium text-xl'>To complete registration, follow the link sent to the email address you provided</div>
                                <div className='mt-6'>
                                    <Link to='/auth/signin' className='text-primary'>
                                        Sign in
                                    </Link>
                                </div>
                            </>
                        )}
                        {!accountCreated &&
                            <form onSubmit={handleSubmit}>
                                <div className='mb-4'>
                                    <label className='mb-2.5 block font-medium text-black dark:text-white'>
                                        Email
                                    </label>
                                    <TextInput
                                        rightIcon={HiMiniEnvelope}
                                        sizing='lg'
                                        type='email'
                                        placeholder='Enter your email'
                                        value={email}
                                        onChange={(event: React.ChangeEvent<HTMLInputElement>) => setEmail(event.target.value)}
                                        required
                                        theme={inputTheme}
                                    />
                                </div>
                                <div className='mb-4'>
                                    <label className='mb-2.5 block font-medium text-black dark:text-white'>
                                        Password
                                    </label>
                                    <TextInput
                                        rightIcon={HiMiniLockClosed}
                                        sizing='lg'
                                        type='password'
                                        placeholder='Enter your password'
                                        value={password}
                                        onChange={handlePasswordChange}
                                        helperText={passwordError || PASSWORD_HELPER_TEXT}
                                        color={passwordError ? 'failure' : undefined}
                                        required
                                        theme={inputTheme}
                                    />
                                </div>
                                <div className='mb-6'>
                                    <label className='mb-2.5 block font-medium text-black dark:text-white'>
                                        Repeate Password
                                    </label>
                                    <TextInput
                                        rightIcon={HiMiniLockClosed}
                                        sizing='lg'
                                        type='password'
                                        placeholder='Repeate your password'
                                        value={repeatePassword}
                                        onChange={handleRepeatePasswordChange}
                                        helperText={repeatePasswordError}
                                        color={repeatePasswordError ? 'failure' : undefined}
                                        required
                                        theme={inputTheme}
                                    />
                                </div>
                                <div className="mb-5">
                                    <Button size='lg' className='w-full' disabled={submitDisabled} type='submit' color='blue'>
                                        {loading ? <Spinner /> : 'Create account'}
                                    </Button>
                                </div>
                                <div className='mt-6 text-center'>
                                    <p>
                                        Already have an account?{' '}
                                        <Link to='/auth/signin' className='text-primary'>
                                            Sign in
                                        </Link>
                                    </p>
                                </div>
                            </form>
                        }
                    </div>
                </div>
            </div>
        </div>
    );
};

export default SignUp;
