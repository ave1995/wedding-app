import { useEffect, useState } from "react";
import { apiUrl } from "../../functions/api";
import { get } from "../../functions/fetch";
import { useApiErrorHandler } from "../../hooks/useApiErrorHandler";
import { PieChart, pieArcLabelClasses } from "@mui/x-charts/PieChart";
import { legendClasses } from "@mui/x-charts/ChartsLegend";

export type QuestionStats = {
  right: number;
  wrong: number;
};

interface QuestionStatsComp {
  id: string;
}

export default function QuestionStats({ id }: QuestionStatsComp) {
  const { handleError } = useApiErrorHandler();
  const [stats, setStats] = useState<QuestionStats>({ right: 0, wrong: 0 });

  useEffect(() => {
    const fetchStats = async () => {
      const result = await get<QuestionStats>(
        apiUrl(`/api/questions/${id}/stats`),
        null,
        true
      );
      if (handleError(result.error, result.status)) return;

      setStats(result.data!);
    };

    fetchStats();
  }, [id]);

  if (!stats) {
    return <p>Načítám statistiky...</p>;
  }

  return (
    <PieChart
      series={[
        {
          data: [
            {
              value: stats.wrong,
              label: "Špatně",
              labelMarkType: "circle",
              color: "#f6339a",
            },
            {
              value: stats.right,
              label: "Správně",
              labelMarkType: "circle",
              color: "#38C172",
            },
          ],
          arcLabel: (item) => `${item.value}`,
        },
      ]}
      sx={{
        [`& .${pieArcLabelClasses.root}`]: {
          fill: "white",
          fontSize: 16,
          fontWeight: 600,
        },
        [`.${legendClasses.label}`]: {
          fontSize: 16,
          fontWeight: 600,
        },
        [`.${legendClasses.mark}`]: {
          height: 18,
          width: 18,
        },
      }}
      width={200}
      height={200}
    />
  );
}
