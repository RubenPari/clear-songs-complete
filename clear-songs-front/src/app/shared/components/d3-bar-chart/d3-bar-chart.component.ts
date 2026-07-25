/**
 * D3 bar chart component for visualizing artist track counts.
 * Features responsive design, smooth animations, and interactive tooltips.
 */
import { Component, Input, OnChanges, SimpleChanges, ElementRef, ViewChild, AfterViewInit, OnDestroy, HostListener } from '@angular/core';
import * as d3 from 'd3';

/* eslint-disable @typescript-eslint/no-explicit-any */

/** Data point for the bar chart. */
interface ChartData {
  label: string;
  value: number;
}

/**
 * Interactive bar chart component using D3.js.
 * Renders a vertical bar chart with animations and hover tooltips.
 */
@Component({
  selector: 'app-d3-bar-chart',
  templateUrl: './d3-bar-chart.component.html',
  styleUrls: ['./d3-bar-chart.component.scss'],
  standalone: true
})
export class D3BarChartComponent implements OnChanges, AfterViewInit, OnDestroy {
  @ViewChild('chartContainer', { static: true }) chartContainer!: ElementRef<HTMLDivElement>;

  @Input() data: ChartData[] = [];
  @Input() height = 300;
  /** Theme-reactive colors using CSS custom properties. */
  @Input() colors: string[] = [
    'hsl(var(--primary))',
    'hsl(var(--primary) / 0.85)',
    'hsl(var(--primary) / 0.7)',
    'hsl(var(--primary) / 0.55)',
    'hsl(var(--primary) / 0.42)'
  ];

  private svg: any;
  private margin = { top: 20, right: 20, bottom: 60, left: 50 };
  private width = 0;
  private chartHeight = 0;
  private xScale: any;
  private yScale: any;
  private tooltip: any;

  ngAfterViewInit(): void {
    this.initChart();
    if (this.data && this.data.length > 0) {
      this.updateChart();
    }
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['data'] && !changes['data'].firstChange && this.svg) {
      if (this.data && this.data.length > 0) {
        this.updateChart();
      } else {
        this.svg.selectAll('.bar').remove();
      }
    }
  }

  ngOnDestroy(): void {
    if (this.tooltip) {
      this.tooltip.remove();
    }
  }

  /** Initializes the SVG container, scales, axes, and tooltip. */
  private initChart(): void {
    const container = this.chartContainer.nativeElement;

    if (container.clientWidth === 0) {
      setTimeout(() => this.initChart(), 100);
      return;
    }

    this.width = container.clientWidth - this.margin.left - this.margin.right;
    this.chartHeight = this.height - this.margin.top - this.margin.bottom;

    d3.select(container).select('svg').remove();

    this.svg = d3.select(container)
      .append('svg')
      .attr('width', this.width + this.margin.left + this.margin.right)
      .attr('height', this.height + 50)
      .append('g')
      .attr('transform', `translate(${this.margin.left},${this.margin.top})`);

    this.tooltip = d3.select('body')
      .append('div')
      .attr('class', 'd3-chart-tooltip')
      .style('opacity', 0)
      .style('position', 'absolute')
      .style('background', 'hsl(var(--popover))')
      .style('color', 'hsl(var(--popover-foreground))')
      .style('border', '1px solid hsl(var(--border))')
      .style('padding', '10px 14px')
      .style('border-radius', '10px')
      .style('font-size', '13px')
      .style('font-weight', '600')
      .style('font-family', 'var(--font-display)')
      .style('pointer-events', 'none')
      .style('z-index', '10000')
      .style('box-shadow', '0 10px 30px -10px rgba(0, 0, 0, 0.35)');

    this.xScale = d3.scaleBand()
      .range([0, this.width])
      .padding(0.3);

    this.yScale = d3.scaleLinear()
      .range([this.chartHeight, 0]);

    this.svg.append('g')
      .attr('class', 'x-axis')
      .attr('transform', `translate(0,${this.chartHeight})`);

    this.svg.append('g')
      .attr('class', 'y-axis');

    this.svg.append('g')
      .attr('class', 'grid-lines');
  }

  /** Updates the chart with new data, animating bars and axes. */
  private updateChart(): void {
    if (!this.data || this.data.length === 0) {
      return;
    }

    this.xScale.domain(this.data.map(d => d.label));
    const maxValue = d3.max(this.data, d => d.value) || 0;
    this.yScale.domain([0, maxValue + Math.ceil(maxValue * 0.1)]);

    this.svg.select('.x-axis')
      .transition()
      .duration(500)
      .call(d3.axisBottom(this.xScale))
      .selectAll('text')
      .style('text-anchor', 'end')
      .attr('dx', '-.8em')
      .attr('dy', '.15em')
      .attr('transform', 'rotate(-45)')
      .style('font-size', '12px')
      .style('font-weight', '600')
      .style('font-family', 'var(--font-display)')
      .style('fill', 'hsl(var(--muted-foreground))');

    this.svg.select('.y-axis')
      .transition()
      .duration(500)
      .call(
        d3.axisLeft(this.yScale)
          .ticks(Math.min(maxValue, 10))
          .tickFormat(d => d.toString())
      )
      .selectAll('text')
      .style('font-size', '12px')
      .style('font-weight', '600')
      .style('font-family', 'var(--font-display)')
      .style('fill', 'hsl(var(--muted-foreground))');

    // Update grid lines using D3 enter/update/exit pattern.
    this.svg.select('.grid-lines')
      .selectAll('line')
      .data(this.yScale.ticks(Math.min(maxValue, 10)))
      .join(
        (enter: any) => enter.append('line')
          .attr('class', 'grid-line')
          .attr('x1', 0)
          .attr('x2', this.width)
          .attr('y1', (d: number) => this.yScale(d))
          .attr('y2', (d: number) => this.yScale(d))
          .style('stroke', 'hsl(var(--border))')
          .style('stroke-width', 1)
          .style('stroke-dasharray', '3,3'),
        (update: any) => update
          .transition()
          .duration(500)
          .attr('y1', (d: number) => this.yScale(d))
          .attr('y2', (d: number) => this.yScale(d)),
        (exit: any) => exit.remove()
      );

    this.svg.selectAll('.bar').remove();

    const bars = this.svg.selectAll('.bar')
      .data(this.data)
      .enter()
      .append('rect')
      .attr('class', 'bar')
      .attr('x', (d: ChartData) => this.xScale(d.label)!)
      .attr('width', this.xScale.bandwidth())
      .attr('y', this.chartHeight)
      .attr('height', 0)
      .attr('fill', (d: ChartData, i: number) => this.colors[i % this.colors.length])
      .attr('rx', 8)
      .attr('ry', 8)
      .style('cursor', 'pointer')
      .style('transition', 'all 0.2s ease');

    bars.transition()
      .duration(800)
      .ease(d3.easeCubicOut)
      .attr('y', (d: ChartData) => this.yScale(d.value))
      .attr('height', (d: ChartData) => this.chartHeight - this.yScale(d.value));

    // Add hover effects with tooltip positioning.
    bars.on('mouseover', (event: MouseEvent, d: ChartData) => {
        d3.select(event.currentTarget as any)
          .transition()
          .duration(200)
          .attr('opacity', 0.9)
          .attr('transform', 'scale(1.02)');

        this.tooltip
          .style('opacity', 1)
          .html(`<div style="margin-bottom: 4px; font-size: 14px;">${d.label}</div><div style="color: hsl(var(--primary));">${d.value} tracks</div>`)
          .style('left', (event.pageX + 10) + 'px')
          .style('top', (event.pageY - 10) + 'px');
      })
      .on('mouseout', (event: MouseEvent) => {
        d3.select(event.currentTarget as any)
          .transition()
          .duration(200)
          .attr('opacity', 1)
          .attr('transform', 'scale(1)');

        this.tooltip.style('opacity', 0);
      });

    this.svg.selectAll('.x-axis line, .y-axis line')
      .style('stroke', 'transparent');

    this.svg.selectAll('.x-axis path, .y-axis path')
      .style('stroke', 'transparent');
  }

  @HostListener('window:resize')
  onResize(): void {
    if (this.chartContainer && this.svg) {
      const container = this.chartContainer.nativeElement;
      this.width = container.clientWidth - this.margin.left - this.margin.right;

      d3.select(container).select('svg')
        .attr('width', this.width + this.margin.left + this.margin.right);

      this.xScale.range([0, this.width]);

      this.svg.selectAll('.grid-line')
        .attr('x2', this.width);

      this.svg.selectAll('.bar')
        .attr('x', (d: ChartData) => this.xScale(d.label)!)
        .attr('width', this.xScale.bandwidth());

      this.svg.select('.x-axis')
        .attr('transform', `translate(0,${this.chartHeight})`)
        .call(d3.axisBottom(this.xScale));
    }
  }
}
