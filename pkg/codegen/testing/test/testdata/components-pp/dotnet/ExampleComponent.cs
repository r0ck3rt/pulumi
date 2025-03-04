using System.Collections.Generic;
using Pulumi;
using Random = Pulumi.Random;

namespace Components
{
    public class ExampleComponentArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// A simple input
        /// </summary>
        [Input("input")]
        public Input<string> Input { get; set; } = null!;
        /// <summary>
        /// The main CIDR blocks for the VPC
        /// It is a map of strings
        /// </summary>
        [Input("cidrBlocks")]
        public InputMap<string> CidrBlocks { get; set; } = null!;
        [Input("ipAddress")]
        public InputList<int> IpAddress { get; set; } = null!;
    }

    public class ExampleComponent : global::Pulumi.ComponentResource
    {
        [Output("result")]
        public Output<string> Result { get; private set; }
        public ExampleComponent(string name, ExampleComponentArgs args, ComponentResourceOptions? opts = null)
            : base("components:index:ExampleComponent", name, args, opts)
        {
            var password = new Random.RandomPassword($"{name}-password", new()
            {
                Length = 16,
                Special = true,
                OverrideSpecial = args.Input,
            }, new CustomResourceOptions
            {
                Parent = this,
            });

            var simpleComponent = new Components.SimpleComponent($"{name}-simpleComponent", new ComponentResourceOptions
            {
                Parent = this,
            });

            this.Result = password.Result;

            this.RegisterOutputs(new Dictionary<string, object?>
            {
                ["result"] = password.Result,
            });
        }
    }
}
