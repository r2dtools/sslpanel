import { Button } from 'flowbite-react';
import React from 'react';
import { HiMiniArrowLeft } from "react-icons/hi2";
import Breadcrumb from '../components/Breadcrumb';
import image from '../images/403.svg';

const Error403: React.FC = () => {
    return (
        <div>
            <Breadcrumb pageName='Error Page' />
            <div className="rounded-sm border border-stroke bg-white px-5 py-10 shadow-default dark:border-strokedark dark:bg-boxdark sm:py-20">
                <div className="mx-auto max-w-[410px]">
                    <img src={image} />
                    <div className="mt-7.5 text-center">
                        <h2 className="mb-3 text-2xl font-bold text-black dark:text-white">
                            Access denied
                        </h2>
                        <p className="font-medium">
                            Sorry, you do not have access to this page.
                        </p>
                        <Button href='/' color='blue' className='mt-7.5 inline-flex'><HiMiniArrowLeft className="mr-2 h-5 w-5" />Back to Home</Button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Error403;
