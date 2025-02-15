import Breadcrumb from '../../components/Breadcrumb';
import { Button } from 'flowbite-react';
import CertidicateItem from '../../features/certificate/components/CertificateItem';
import { AMAZON_CODE, CLOUD_FLARE_CODE, COMODO_CODE, DIGICERT_CODE, LE_CODE, RAPID_SSL_CODE, SECTIGO_CODE } from '../../features/certificate/constants';

const CertificateList = () => {
    return (
        <>
            <Breadcrumb pageName='Certificates'>
                <Button color='blue'>Add</Button>
            </Breadcrumb>
            <div className="rounded-sm border border-stroke bg-white shadow-default dark:border-strokedark dark:bg-boxdark">
                <div className="border-b border-stroke px-4 py-5 dark:border-strokedark md:px-6 xl:px-7.5">
                    <div className="flex items-center gap-3">
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
                    <div className="flex flex-col gap-7">
                        <CertidicateItem code={LE_CODE} />
                        <CertidicateItem code={SECTIGO_CODE} />
                        <CertidicateItem code={DIGICERT_CODE} />
                        <CertidicateItem code={CLOUD_FLARE_CODE} />
                        <CertidicateItem code={COMODO_CODE} />
                        <CertidicateItem code={AMAZON_CODE} />
                        <CertidicateItem />
                        <CertidicateItem code={RAPID_SSL_CODE} />
                    </div>
                </div>
            </div>
        </>
    );
}

export default CertificateList;
