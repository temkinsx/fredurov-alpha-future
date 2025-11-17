import {FileText, Scale, TrendingUpIcon} from "lucide-react";
import ActionCard from "./ActionCard.tsx";


function QuickActions() {


    return (
        <div className="flex flex-col gap-4">
            <h2>Быстрые действия</h2>
            <div className="flex gap-5">
                <ActionCard
                    Icon={ TrendingUpIcon }
                    header={"Финансовая аналитика"}
                    details={"Аналитика финансовых показателей"}
                />
                <ActionCard
                    Icon={ FileText }
                    header={"Шаблоны маркетинга"}
                    details={"Готовые маркетинг-решения"}
                />
                <ActionCard
                    Icon={ Scale }
                    header={"Правовые документы"}
                    details={"Юридическая поддержка"}
                />
            </div>
        </div>
    );
}

export default QuickActions;
