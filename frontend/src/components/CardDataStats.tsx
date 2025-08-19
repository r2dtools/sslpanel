import React, { ReactNode } from 'react';
import { HiOutlineArrowRightCircle } from 'react-icons/hi2';

interface CardDataStatsProps {
    title: string;
    total: number;
    children: ReactNode;
}

const CardDataStats: React.FC<CardDataStatsProps> = ({ title, total, children }) => {
    return (
        <div className="rounded-sm border border-stroke bg-white py-6 px-7.5 shadow-default dark:border-strokedark dark:bg-boxdark">
            <div className="flex h-11.5 w-11.5 items-center justify-center rounded-full bg-meta-2 dark:bg-meta-4">
                {children}
            </div>
            <div className="mt-4 flex items-end justify-between">
                <div>
                    <h4 className="text-xl font-bold text-black dark:text-white">
                        {title}
                    </h4>
                    <span className="text-sm font-medium">{`Total ${total}`}</span>
                </div>
                <HiOutlineArrowRightCircle size={24} />
            </div>
        </div>
    );
};

export default CardDataStats;
