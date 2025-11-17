import type { ComponentType } from "react";
import type { LucideProps } from "lucide-react";


interface ActionCardProps {
    Icon: ComponentType<LucideProps>;
    header: string;
    details: string;
}

function ActionCard( { Icon, header, details }: ActionCardProps) {

    const handleQuickAction = () => {

    }

    return (
        <div className="flex flex-col gap-3 p-3 justify-center
        bg-white max-w-50 min-h-50 rounded-2xl border border-gray-200
        hover:bg-gray-50 transition-all shadow-xl shadow-gray-200"
            onClick={ handleQuickAction }
        >
            <Icon className="p-3 h-13 w-13 bg-gray-100 border border-gray-300 text-black rounded-2xl"/>
            <h3>{ header }</h3>
            <p className="text-gray-400 text-sm">{ details }</p>
        </div>
    );
}

export default ActionCard;