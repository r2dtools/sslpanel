import React, { useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Spinner } from 'flowbite-react';
import { toast } from 'react-toastify';
import { confirm } from '../../features/auth/authApi';
import AuthSide from '../../components/AuthSide';

const Confirm: React.FC = () => {
    const { search } = useLocation();
    const navigate = useNavigate();

    const query = new URLSearchParams(search);
    const token = query.get('token') || '';

    useEffect(() => {
        confirmEmail(token);
    }, [token]);

    const confirmEmail = async (token: string) => {
        if (!token) {
            toast.error('Invalid token');
            navigate('/auth/sigin');

            return;
        }

        try {
            await confirm(token);
            toast.success('Email successfully verified!');
        } catch (error) {
            toast.error((error as Error).message);
        } finally {
            navigate('/auth/sigin');
        }
    };

    return (
        <>
            <div className="rounded-sm border border-stroke bg-white shadow-default dark:border-strokedark dark:bg-boxdark">
                <div className="flex flex-wrap items-center">
                    <div className="hidden w-full xl:block xl:w-1/2">
                        <AuthSide />
                    </div>
                    <div className="w-full border-stroke dark:border-strokedark xl:w-1/2 xl:border-l-2">
                        <div className="w-full p-4 sm:p-12.5 xl:p-17.5">
                            <h2 className="mb-9 text-2xl font-bold text-black dark:text-white sm:text-title-xl2">
                                Email confirmation
                            </h2>
                            <div className='flex gap-2'>
                                <div className='font-medium text-xl'>Please wait ...</div>
                                <Spinner />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </>
    );
};

export default Confirm;
