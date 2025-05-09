import { Button, HRTrimmed, Label, Spinner, TextInput } from 'flowbite-react';
import Breadcrumb from '../components/Breadcrumb';
import { useEffect, useState } from 'react';
import { useAppDispatch, useAppSelector } from '../app/hooks';
import { selectCurrentUser } from '../features/auth/authSlice';
import { HiOutlineEnvelope } from 'react-icons/hi2';
import { createPasswordValidator, PASSWORD_HELPER_TEXT } from '../lib/form';
import useAuthToken from '../features/auth/hooks';
import { changePassword, selectPasswordChangeStatus } from '../features/account/accountSlice';
import Error403 from './Error403';
import { FetchStatus } from '../app/types';

const Settings = () => {
    const [authToken] = useAuthToken();
    const dispatch = useAppDispatch();
    const currentUser = useAppSelector(selectCurrentUser);
    const passwordChangeStatus = useAppSelector(selectPasswordChangeStatus);
    const [currentPassword, setCurrentPassword] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [repeatePassword, setRepeatePassword] = useState<string>('');
    const [passwordError, setPasswordError] = useState<string>('');
    const [repeatePasswordError, setRepeatePasswordError] = useState<string>('');

    if (!authToken) {
        return <Error403 />
    }

    useEffect(() => {
        if (passwordChangeStatus === FetchStatus.Succeeded) {
            resetPasswordForm();
        }
    }, [passwordChangeStatus]);

    const handlePasswordChange = async (event: React.FormEvent<HTMLFormElement>) => {
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

        await dispatch(changePassword({
            password: currentPassword,
            newPassword: password,
            token: authToken,
        }));
    };

    const resetPasswordForm = () => {
        setCurrentPassword('');
        setPassword('');
        setRepeatePassword('');
    };

    const handleNewPasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setPasswordError('');
        setPassword(event.target.value);
    };

    const handleRepeatePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRepeatePasswordError('');
        setRepeatePassword(event.target.value);
    };

    return (
        <>
            <div className="mx-auto max-w-3xl">
                <Breadcrumb pageName="Settings" />

                <div className="mt-3">
                    <Label htmlFor="email" className="mb-2 block">
                        Email
                    </Label>
                    <TextInput
                        id="email"
                        name="email"
                        value={currentUser?.email}
                        icon={HiOutlineEnvelope}
                        disabled
                    />
                </div>
                <HRTrimmed className='bg-gray-200 dark:bg-gray-700' />

                <h3 className='text-lg font-medium mb-6'>Change Password</h3>
                <form onSubmit={handlePasswordChange}>
                    <div className="mb-6">
                        <Label htmlFor="currentPassword" className="mb-2 block">
                            Current Password
                        </Label>
                        <TextInput
                            id="currentPassword"
                            name="currentPassword"
                            placeholder="Please enter your current password"
                            value={currentPassword}
                            onChange={(event: React.ChangeEvent<HTMLInputElement>) => setCurrentPassword(event.target.value)}
                            type='password'
                            required
                        />
                    </div>
                    <div className="mb-6">
                        <Label htmlFor="password" className="mb-2 block">
                            New Password
                        </Label>
                        <TextInput
                            id="password"
                            name="password"
                            placeholder="Please enter a new password"
                            value={password}
                            onChange={handleNewPasswordChange}
                            helperText={<>{passwordError || PASSWORD_HELPER_TEXT}</>}
                            color={passwordError ? 'failure' : undefined}
                            type='password'
                            required
                        />
                    </div>
                    <div className="mb-6">
                        <Label htmlFor="repeatePassword" className="mb-2 block">
                            Repeate Password
                        </Label>
                        <TextInput
                            id="repeatePassword"
                            name="repeatePassword"
                            placeholder="Please repeate password"
                            value={repeatePassword}
                            onChange={handleRepeatePasswordChange}
                            helperText={<>{repeatePasswordError}</>}
                            color={repeatePasswordError ? 'failure' : undefined}
                            type='password'
                            required
                        />
                    </div>
                    <Button className="w-full" color='blue' type='submit'>
                        {passwordChangeStatus === FetchStatus.Pending ? <Spinner /> : 'Change'}
                    </Button>
                </form>
            </div>
        </>
    );
};

export default Settings;
