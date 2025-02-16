import Breadcrumb from '../../components/Breadcrumb';
import { Button } from 'flowbite-react';
import CertificateItem from '../../features/certificate/components/CertificateItem';
import useAuthToken from '../../features/auth/hooks';
import { useAppDispatch, useAppSelector } from '../../app/hooks';
import { fetchCertificates, selectCertificates, selectCertificatesFetchStatus } from '../../features/certificate/certificatesSlice';
import { useEffect } from 'react';
import { FetchStatus } from '../../app/types';
import { useParams } from 'react-router-dom';
import Loader from '../../components/Loader/Loader';
import { Certificate } from '../../features/certificate/types';
import empty from '../../images/empty.png';

const CertificateList = () => {
    const { guid } = useParams();
    const [authToken] = useAuthToken();
    const dispatch = useAppDispatch();
    const certificatesSelectStatus = useAppSelector(selectCertificatesFetchStatus);
    const certificates = useAppSelector(selectCertificates);

    const isLoading = certificatesSelectStatus === FetchStatus.Pending;
    const certs = Object.values(certificates);

    useEffect(() => {
        if (authToken && certificatesSelectStatus !== FetchStatus.Pending) {
            dispatch(fetchCertificates({
                guid: guid as string,
                token: authToken,
            }));
        }
    }, [authToken]);

    return (
        isLoading
            ? <Loader />
            : <>
                <Breadcrumb pageName='Certificates'>
                    <Button color='blue'>Add</Button>
                </Breadcrumb>
                <div className="rounded-sm border border-stroke bg-white shadow-default dark:border-strokedark dark:bg-boxdark">
                    <div className="border-b border-stroke px-4 py-5 dark:border-strokedark md:px-6 xl:px-7.5">
                        <div className="flex items-center gap-2">
                            <div className="w-3/12 md:w-4/12">
                                <span className="font-medium">Vendor</span>
                            </div>
                            <div className="w-4/12">
                                <span className="font-medium">Names</span>
                            </div>
                            <div className="w-3/12">
                                <span className="font-medium">Expires</span>
                            </div>
                            <div className="w-2/12 md:w-1/12 text-center">
                                <span className="font-medium">Actions</span>
                            </div>
                        </div>
                    </div>
                    <div className="p-4 md:p-6 xl:p-7.5">
                        <div className="mx-auto max-w-[80px]">
                            {!certs.length && <img src={empty} />}
                        </div>
                        {
                            certs.map((certificate: Certificate) => (
                                <CertificateItem certificate={certificate} key={certificate.cn} />
                            ))
                        }
                    </div>
                </div>
            </>
    );
}

export default CertificateList;
