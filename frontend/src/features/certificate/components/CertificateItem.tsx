import { Link } from 'react-router-dom';
import { Badge, Tooltip } from 'flowbite-react';
import { HiMiniTrash } from 'react-icons/hi2';
import { getCertificateIssuerIcon } from '../utils';
import sslIcon from '../../../images/certificate/ca.png';

interface CertificateItemProps {
    code?: string;
};

const CertidicateItem: React.FC<CertificateItemProps> = ({ code }) => {
    const handleDelete = async (event: React.MouseEvent<SVGElement, MouseEvent>) => {
        event.preventDefault();
        event.stopPropagation();

        console.log("delete");
    };

    const preventClick = (event: React.MouseEvent) => {
        event.preventDefault();
        event.stopPropagation();
    };

    const icon = getCertificateIssuerIcon(code);

    return (
        <Link to="#">
            <div className="p-3 flex items-center gap-3 hover:bg-[#F8FAFD] dark:hover:bg-meta-4 hover:rounded">
                <div className="w-3/12 md:w-4/12">
                    {
                        icon ? (
                            <div className="max-w-40">
                                <img src={icon} className='w-full' />
                            </div>
                        ) : (
                            <div className='flex items-center gap-4'>
                                <div className="2xsm:h-11 2xsm:max-w-11">
                                    <img src={sslIcon} />
                                </div>
                                <div>
                                    <span className="font-bold">Default CA</span>
                                </div>
                            </div>
                        )
                    }

                </div>
                <div className="w-4/12 flex flex-col gap-1 font-medium">
                    <div className='truncate'>r2dtools.work.gd</div>
                    <div className='truncate'>www.r2dtools.work.gd</div>
                </div>
                <div className="w-3/12 flex flex-col gap-1 lg:gap-2 items-center lg:flex-row">
                    <span className="font-medium">23 oct 2025</span>
                    <Badge color='success' className='inline'>30 days</Badge>
                </div>
                <div className="w-2/12 text-center lg:w-1/12" onClick={preventClick}>
                    <button className="flex justify-between mx-auto block ">
                        <Tooltip content="Delete" >
                            <HiMiniTrash size={20} className='hover:text-red-500' onClick={handleDelete} />
                        </Tooltip>
                    </button>
                </div>
            </div>
        </Link>
    );
};

export default CertidicateItem;
