import { t } from "@lingui/core/macro"
import type { ChartData } from "@/types"
import { ChartCard } from "../chart-card"
import LineChartDefault from "@/components/charts/line-chart"
import { decimalString } from "@/lib/utils"

/** Open TCP connections on the ports the agent counts (CONNECTION_PORTS):
 * shown only for systems whose agent reports them. */
export function ConnectionsChart({
	chartData,
	grid,
	dataEmpty,
}: {
	chartData: ChartData
	grid: boolean
	dataEmpty: boolean
}) {
	if (!chartData.systemStats.some((record) => record.stats?.conn !== undefined)) {
		return null
	}
	return (
		<ChartCard
			empty={dataEmpty}
			grid={grid}
			title={t`Connections`}
			description={t`Open TCP connections on the agent's CONNECTION_PORTS`}
		>
			<LineChartDefault
				chartData={chartData}
				contentFormatter={(item) => decimalString(item.value, 0)}
				tickFormatter={(value) => String(Math.round(value))}
				dataPoints={[
					{
						label: t`Connections`,
						color: "hsl(160, 84%, 39%)", // Green
						dataKey: ({ stats }) => stats?.conn ?? 0,
					},
				]}
			></LineChartDefault>
		</ChartCard>
	)
}
