import React, { ReactNode } from 'react';
import { HiOutlineArrowRightCircle } from 'react-icons/hi2';
import { Link } from 'react-router-dom';

interface CardStatsLink {
    title: string;
    total: number;
    url: string;
    children: ReactNode;
}

const CardStatsLink: React.FC<CardStatsLink> = ({ title, total, url, children }) => {
    return (
        <Link to={url}>
            <div className="rounded-sm border border-stroke bg-white py-6 px-7.5 shadow-default dark:border-strokedark dark:bg-boxdark hover:bg-[#F8FAFD] dark:hover:bg-meta-4">
                <div className="flex h-11.5 w-11.5 items-center justify-center rounded-full bg-meta-2 dark:bg-meta-4 border-2 border-solid">
                    {children}
                </div>
                <div className="mt-4 flex items-end justify-between">
                    <div>
                        <h4 className="text-title-md font-bold text-black dark:text-white">
                            {title}
                        </h4>
                        <span className="text-sm font-medium">{`Total ${total}`}</span>
                    </div>
                    <HiOutlineArrowRightCircle size={24} />
                </div>
            </div>
        </Link>
    );
};

export default CardStatsLink;
